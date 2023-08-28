

# Copy the contents of the source folder to the destination folder
Copy-Item -Path "../pkg/wubi/data" -Destination "." -Recurse -Force
# Copy the file to the destination folder
Copy-Item -Path "../pinyin-data/pinyin.txt" -Destination "data"
Copy-Item -Path "../pinyin-data/duoyin.txt" -Destination "data"
Copy-Item -Path "../pinyin-data/correct.txt" -Destination "data"
Copy-Item -Path "../UnicodeCJK-WuBi/CJK.txt" -Destination "data"

cd ..

Write-Output "编译 windows 版本"
go build -o build/rose.exe

cd build
