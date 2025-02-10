./scripts/download_jsonl/download.py exp-models/GAIR-LIMO --output datasets/GAIR-LIMO.jsonl && \
    go run ./scripts/translate_dataset \
        -in datasets/GAIR-LIMO.jsonl \
        -out datasets/GAIR-LIMO-KOREAN.jsonl \
        -src English \
        -dst Korean