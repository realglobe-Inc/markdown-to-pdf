#!/bin/sh -e

# MarkdownファイルをPDFに変換する
# 引数でファイルパスを1つ受け取る
# 受け取ったファイルが.mdならPDFに変換してPDFファイルのパスを表示する
# 受け取ったファイルが.tarなら
#   展開して.mdファイルをPDFに変換する
#   できたPDFファイルが1つならPDFファイルのパスを表示する
#   できたPDFファイルが複数ならPDFファイルをTARで固めてTARファイルのパスを表示する

if ! [ -e "${1}" ]; then
  echo "source file does not exist" >&2
  exit 1
fi

work_dir="$(dirname "${1}")"
filename="$(basename "${1}")"

(
  cd "${work_dir}"
  work_dir="$(pwd)"

  case "${filename}" in
  *.md)
    # 何もしない
    ;;
  *.tar)
    tar xf "${filename}"
    ;;
  esac

  md_list=.markdown.list
  find . -path '*.md' -printf '%P\n' >"${md_list}"

  pdf_dir=.pdf
  rm -rf "${pdf_dir}"
  mkdir "${pdf_dir}"
  pdf_list=.pdf.list
  : >"${pdf_list}"
  while read -r md_file; do
    pdf_file="${pdf_dir}/${md_file%.md}.pdf"
    mkdir -p "$(dirname "${pdf_file}")"
    pandoc --pdf-engine=xelatex -V documentclass=bxjsarticle -V classoption=pandoc --self-contained --resource-path=. "${md_file}" -o "${pdf_file}"
    echo "${pdf_file}" >>"${pdf_list}"
  done <"${md_list}"

  count="$(wc -l <"${pdf_list}")"
  if [ "${count}" = 1 ]; then
    echo "${work_dir}/$(head -n 1 "${pdf_list}")"
  elif [ "${count}" -gt 1 ]; then
    # tar -C "${pdf_dir}" cf "${work_dir}/pdf.tar" *
    tar cf pdf.tar pdf
    echo "${work_dir}/pdf.tar"
  else
    echo "no markdown" >&2
    exit 1
  fi
)
