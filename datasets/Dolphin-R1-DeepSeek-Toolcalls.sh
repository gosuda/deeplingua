./scripts/download_jsonl/download.py exp-models/dolphin-r1-deepseek-toolcalls --output datasets/dolphin-r1-deepseek-toolcalls.jsonl && \
    go run ./scripts/translate_dataset \
        -in datasets/dolphin-r1-deepseek-toolcalls.jsonl \
        -out datasets/dolphin-r1-korean-deepseek-toolcalls.jsonl \
        -src English \
        -dst Korean

go run ./scripts/fix_jsonl datasets/dolphin-r1-korean-deepseek-toolcalls.jsonl datasets/dolphin-r1-korean-deepseek-toolcalls.fixed.jsonl

python3 scripts/upload_jsonl/upload.py \
    --name exp-models/dolphin-r1-korean-deepseek-toolcalls \
    datasets/dolphin-r1-korean-deepseek-toolcalls.fixed.jsonl