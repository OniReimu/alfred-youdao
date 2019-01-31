#!/bin/bash

# 可以指定语言发音
# say -l zn_CN 我爱你中国

PATH='/bin:/usr/bin'
POSITIONAL=()
LANGUAGE=''
VOICE=''

while [[ $# -gt 0 ]]
do
    key=$1

    case $key in
        -l|--language)
            LANGUAGE=$2
            shift
            shift
            ;;
        *)
            POSITIONAL+=("$1")
            shift
            ;;
    esac
done

set -- "${POSITIONAL[@]}" # 数组转化成字符串

if [[ ${#LANGUAGE} -gt 0 ]]; then # ${#var} 相当于 python里的len(var)
    VOICE=$(say -v ? | grep $LANGUAGE | cut -f1 -d ' ' | head -n 1)
fi

if [[ ${#VOICE} -gt 0 ]]; then
    say -v $VOICE $@
else
    say $@
fi
