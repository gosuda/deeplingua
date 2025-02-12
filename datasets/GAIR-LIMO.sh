./scripts/download_jsonl/download.py exp-models/GAIR-LIMO --output datasets/GAIR-LIMO.jsonl && \
    go run ./scripts/translate_dataset \
        -in datasets/GAIR-LIMO.jsonl \
        -out datasets/GAIR-LIMO-KOREAN.jsonl \
        -src English \
        -dst Korean

go run ./scripts/fix_jsonl datasets/GAIR-LIMO-KOREAN.jsonl datasets/GAIR-LIMO-KOREAN.fixed.jsonl

python3 scripts/upload_jsonl/upload.py \
    --name exp-models/GAIR-LIMO-KOREAN \
    datasets/GAIR-LIMO-KOREAN.fixed.jsonl