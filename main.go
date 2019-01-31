package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	al_queue "github.com/AceDarkknight/AlgorithmAndDataStructure/queue"
	"github.com/zgs225/alfred-youdao/alfred"
	"github.com/zgs225/youdao"
)

const (
	APPID     = "2f871f8481e49b4c"
	APPSECRET = "CQFItxl9hPXuQuVcQa5F2iPmZSbN0hYS"
	MAX_LEN   = 255

	UPDATECMD = "alfred-youdao:update"

	QUEUE_SIZE = 20 // Size of the unique queue
)

func init() {
	log.SetPrefix("[i] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	log.Println(os.Args)

	client := &youdao.Client{
		AppID:     APPID,
		AppSecret: APPSECRET,
	}
	agent := newAgent(client)

	// Load the history data from local cache into the queue
	uniqueQueue, _ := al_queue.NewUniqueQueue(QUEUE_SIZE)
	queue, _ := agent.readHistory()
	if len(queue) == 0 {
		log.Println("no history found")
	} else if len(queue) <= uniqueQueue.Capacity() {
		log.Println("history found")
		for _, v := range queue {
			uniqueQueue.Enqueue(v)
		}
	}

	items := alfred.NewResult()
	if os.Args[1] == "#" {
		mod := map[string]*alfred.ModElement{
			alfred.Mods_Cmd: &alfred.ModElement{
				Valid:    true,
				Subtitle: "发音",
			},
		}
		for {
			if uniqueQueue.Length() == 0 {
				break
			}
			output := strings.Split(uniqueQueue.PopItem().(string), "=>")
			if len(output) < 2 {
				log.Println("Fail to retrieve the history data")
				break
			}
			mod2 := copyModElementMap(mod)
			mod2[alfred.Mods_Shift] = &alfred.ModElement{
				Valid:    true,
				Arg:      toYoudaoDictUrl(output[0]),
				Subtitle: "回车键打开词典网页",
			}
			if !containChinese(output[0]) {
				items.Append(&alfred.ResultElement{
					Valid:    true,
					Title:    output[0],
					Subtitle: output[1],
					Arg:      output[0],
					Mods:     mod2,
					// QuickLookUrl: toYoudaoDictUrl(output[0]),
				})
			} else {
				items.Append(&alfred.ResultElement{
					Valid:    true,
					Title:    output[0],
					Subtitle: output[1],
					Arg:      eliminateChinese(strings.Split(output[1], ": ")[0]),
					Mods:     mod2,
					// QuickLookUrl: toYoudaoDictUrl(output[0]),
				})
			}
		}
	} else {
		q, from, to, lang := parseArgs(os.Args)
		// items := alfred.NewResult()

		if lang {
			if err := agent.Client.SetFrom(from); err != nil {
				items.Append(&alfred.ResultElement{
					Valid:    true,
					Title:    fmt.Sprintf("错误: 源语言不支持[%s]", from),
					Subtitle: `有道词典`,
				})
				items.End()
			}
			if err := agent.Client.SetTo(to); err != nil {
				items.Append(&alfred.ResultElement{
					Valid:    true,
					Title:    fmt.Sprintf("错误: 目标语言不支持[%s]", to),
					Subtitle: `有道词典`,
				})
				items.End()
			}
		}

		if len(q) == 0 {
			items.Append(&alfred.ResultElement{
				Valid:    true,
				Title:    "有道词典",
				Subtitle: `查看"..."的解释或翻译`,
			})
			items.End()
		}

		if len(q) > 255 {
			items.Append(&alfred.ResultElement{
				Valid:    false,
				Title:    "错误: 最大查询字符数为255",
				Subtitle: q,
			})
			items.End()
		}

		r, err := agent.Query(q)
		if err != nil {
			panic(err)
		}

		mod := map[string]*alfred.ModElement{
			alfred.Mods_Shift: &alfred.ModElement{
				Valid:    true,
				Arg:      toYoudaoDictUrl(q),
				Subtitle: "回车键打开词典网页",
			},
		}

		if r.Basic != nil {
			var flag bool // A flag to indicate the first result
			phonetic := joinPhonetic(r.Basic.Phonetic, r.Basic.UkPhonetic, r.Basic.UsPhonetic)
			for _, title := range r.Basic.Explains {
				mod2 := copyModElementMap(mod)
				mod2[alfred.Mods_Cmd] = &alfred.ModElement{
					Valid:    true,
					Arg:      wordsToSayCmdOption(title, q, r),
					Subtitle: "发音",
				}
				item := alfred.ResultElement{
					Valid:    true,
					Title:    title,
					Subtitle: phonetic,
					Arg:      title,
					Mods:     mod2,
					// QuickLookUrl: toYoudaoDictUrl(q),
				}
				items.Append(&item)

				// Reshuffle the order when a cache being hit, otherwise directly enqueue.
				if !flag {
					err := uniqueQueue.Enqueue(fmt.Sprintf("%s=>%s", q, title+": "+phonetic))
					if err == al_queue.ErrReplica {
						log.Println("found the cache, reshuffle the order")
						temp_buf := make([]string, 0, QUEUE_SIZE)
						for {
							if temp := uniqueQueue.PopItem(); temp == nil {
								log.Println("pop fails")
								break
							} else if temp.(string) != fmt.Sprintf("%s=>%s", q, title+": "+phonetic) {
								temp_buf = append(temp_buf, temp.(string))
								temp_buf = reverseSlice(temp_buf)
							} else {
								for _, v := range temp_buf {
									uniqueQueue.Enqueue(v)
								}
								uniqueQueue.Enqueue(temp.(string))

								queue = make([]string, 0, QUEUE_SIZE)
								for uniqueQueue.Length() > 0 {
									queue = append(queue, uniqueQueue.Dequeue().(string))
								}
								if err := agent.writeHistory(queue); err != nil {
									log.Println("writeHistory fails")
								} else {
									log.Println("all good")
								}
								break
							}
						}
					} else {
						log.Println("no cache found")
						if uniqueQueue.Length() >= uniqueQueue.Capacity() {
							uniqueQueue.Dequeue()
						}

						queue = make([]string, 0, QUEUE_SIZE)
						for uniqueQueue.Length() > 0 {
							queue = append(queue, uniqueQueue.Dequeue().(string))
						}
						if err := agent.writeHistory(queue); err != nil {
							log.Println("writeHistory fails")
						} else {
							log.Println("all good")
						}
					}
					flag = true
				}
			}
		}

		if r.Translation != nil {
			title := strings.Join(*r.Translation, "; ")
			mod2 := copyModElementMap(mod)
			mod2[alfred.Mods_Cmd] = &alfred.ModElement{
				Valid:    true,
				Arg:      wordsToSayCmdOption(title, q, r),
				Subtitle: "发音",
			}
			item := alfred.ResultElement{
				Valid:    true,
				Title:    title,
				Subtitle: "翻译结果",
				Arg:      title,
				Mods:     mod2,
				// QuickLookUrl: toYoudaoDictUrl(q),
			}
			items.Append(&item)
		}

		if r.Web != nil {
			items.Append(&alfred.ResultElement{
				Valid:    true,
				Title:    "网络释义",
				Subtitle: "有道词典 for Alfred",
			})
			// When from_lan is chinese, the subtitles will be rather output for 网络释义
			for _, elem := range *r.Web {
				mod2 := copyModElementMap(mod)
				mod2[alfred.Mods_Cmd] = &alfred.ModElement{
					Valid: true,
					// Arg:      wordsToSayCmdOption(elem.Key, q, r),
					Arg:      wordsToSayCmdOption(strings.Join(elem.Value, "; "), q, r),
					Subtitle: "发音",
				}
				ls := strings.Split(r.L, "2")
				l_from := languageToSayLanguage(ls[0])
				if l_from == "zh_CN" {
					items.Append(&alfred.ResultElement{
						Valid:    true,
						Title:    elem.Key,
						Subtitle: strings.Join(elem.Value, "; "),
						Arg:      elem.Key,
						Mods:     mod2,
						// QuickLookUrl: toYoudaoDictUrl(q),
					})
				} else {
					items.Append(&alfred.ResultElement{
						Valid:    true,
						Title:    elem.Key,
						Subtitle: strings.Join(elem.Value, "; "),
						Arg:      elem.Key,
						Mods:     mod,
						// QuickLookUrl: toYoudaoDictUrl(q),
					})

				}
				// items.Append(&alfred.ResultElement{
				// 	Valid:    true,
				// 	Title:    elem.Key,
				// 	Subtitle: strings.Join(elem.Value, "; "),
				// 	Arg:      elem.Key,
				// 	Mods:     mod,
				// })
			}
		}
	}

	if agent.Dirty {
		if err := agent.Cache.SaveFile(CACHE_FILE); err != nil {
			log.Println(err)
		}
		// Commit the queue into the local cache file
		if err := agent.History.SaveFile(CACHE_HISTORY_FILE); err != nil {
			log.Println(err)
		}
	}
	items.End()
}
