#!/usr/bin/env python3

import argparse

def main():
    parser = argparse.ArgumentParser(description='Download dataset from Hugging Face Hub and save as JSONL')
    parser.add_argument('path', type=str, help='local .jsonl file path')
    parser.add_argument('--name',type=str, help='Name of the dataset on Hugging Face Hub')
    args = parser.parse_args()

    from datasets import load_dataset
    import os
    import json

    # Set HF_HUB_ENABLE_HF_TRANSFER to 1 to enable faster downloads
    os.environ["HF_HUB_ENABLE_HF_TRANSFER"] = "1"

    # Load the dataset
    dataset = load_dataset("json", data_files=[args.path])
    dataset = dataset.sort("custom_id")
    #dataset =  dataset.remove_columns(["custom_id"])
    print(dataset)
    def map_func(x):
        for i in range(len(x["messages"])):
            role = x["messages"][i]["role"]
            en = x["messages"][i]["content"]
            ko = x["messages"][i]["translated_content"]
            x["messages"][i] = {
                "role": role,
                "content": ko,
                "content_en": en,
            }
        return x
    dataset = dataset.map(map_func)
    print(dataset)
    
    # Push to hub
    dataset.push_to_hub(args.name)

if __name__ == '__main__':
    main()
