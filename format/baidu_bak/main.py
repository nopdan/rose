MASK = 0x2D382324
TABLE = b"qogjOuCRNkfil5p4SQ3LAmxGKZTdesvB6z_YPahMI9t80rJyHW1DEwFbc7nUVX2-"
DECODE_TABLE = bytes.maketrans(TABLE, bytes(range(64)))


print("bytes(range(64)):", bytes(range(64)), "size:", len(bytes(range(64))))
print("TABLE:", TABLE, "size:", len(TABLE))
print("DECODE_TABLE:", DECODE_TABLE, "size:", len(DECODE_TABLE))


for i in DECODE_TABLE:
    print(i)


def decode(data: bytes):
    assert len(data) % 4 == 2
    assert (base64_remainder := data[-2] - 65) in {0, 1, 2} and data[-1] == 0

    data = data.translate(DECODE_TABLE)
    transformed = bytearray()

    for i in range(0, len(data) - 2, 4):
        high_bits = data[i + 3]
        transformed += bytes(
            (
                data[i] | (high_bits & 0b110000) << 2,
                data[i + 1] | (high_bits & 0b1100) << 4,
                data[i + 2] | (high_bits & 0b11) << 6,
            )
        )

    if base64_remainder:
        for i in range(3 - base64_remainder):
            assert transformed.pop() == 0

    print("transformed:", transformed)

    result = bytearray()

    for i in range(0, len(transformed) // 4 * 4, 4):
        chunk = MASK ^ int.from_bytes(transformed[i : i + 4], byteorder="little")
        chunk = (chunk & 0x1FFFFFFF) << 3 | chunk >> 29
        result += chunk.to_bytes(4, byteorder="little")

    print("result:", result)

    if remainder := transformed[i + 4 :]:
        chunk = MASK ^ int.from_bytes(remainder, byteorder="little")
        result += chunk.to_bytes(4, byteorder="little")[: len(remainder)]

    return result


def decode_bin(input_path: str, output_path: str):
    assert input_path.endswith(".bin") and output_path.endswith(".txt")
    with open(input_path, "rb") as f, open(output_path, "wb") as output:
        output.write(b"\xff\xfe")
        f.seek(1)
        for line in f:
            print("line:", line)
            if data := line[1:-1]:
                print("data:", data)
                _data = decode(data)
                output.write(_data)
                # 将 data 按 utf-16 编码转为字符串输出
                print("decode:", _data.decode("utf-16le"))
                output.write(b"\n\0")


if __name__ == "__main__":
    import sys

    decode_bin(
        "../../testdata/百度输入法词库导出_2023_12_28.bin",
        "test_import.txt",
    )

    # decode_bin(sys.argv[1], sys.argv[2])
