./scripts/download_jsonl/download.py exp-models/simplescaling-s1K-1.1 --output datasets/simplescaling-s1K-1.1.jsonl && \
    go run ./scripts/translate_dataset \
        -in datasets/simplescaling-s1K-1.1.jsonl \
        -out datasets/simplescaling-s1K-1.1-korean.jsonl \
        -src English \
        -dst Korean && \
    go run ./scripts/fix_jsonl \
        datasets/simplescaling-s1K-1.1-korean.jsonl \
        datasets/simplescaling-s1K-1.1-korean.fixed.jsonl && \
    python3 scripts/upload_jsonl/upload.py \
        --name lemon-mint/simplescaling-s1K-1.1-korean \
        datasets/simplescaling-s1K-1.1-korean.fixed.jsonl