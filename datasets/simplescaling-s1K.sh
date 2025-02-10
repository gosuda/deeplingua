./scripts/download_jsonl/download.py exp-models/simplescaling-s1K --output datasets/simplescaling-s1K.jsonl && \
    go run ./scripts/translate_dataset \
        -in datasets/simplescaling-s1K.jsonl \
        -out datasets/simplescaling-KS1K.jsonl \
        -src English \
        -dst Korean