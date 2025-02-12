./scripts/download_jsonl/download.py HuggingFaceTB/smol-smoltalk --output datasets/smol-smoltalk.jsonl && \
    go run ./scripts/translate_dataset \
        -in datasets/smol-smoltalk.jsonl \
        -out datasets/smol-koreantalk.jsonl \
        -src English \
        -dst Korean

go run ./scripts/fix_jsonl datasets/smol-koreantalk.jsonl datasets/smol-koreantalk.fixed.jsonl

python3 scripts/upload_jsonl/upload.py \
    --name lemon-mint/smol-koreantalk \
    datasets/smol-koreantalk.fixed.jsonl