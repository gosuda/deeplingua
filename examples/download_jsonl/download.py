#!/usr/bin/env python3

import argparse

def main():
    parser = argparse.ArgumentParser(description='Download dataset from Hugging Face Hub and save as JSONL')
    parser.add_argument('dataset_name', type=str, help='Name of the dataset on Hugging Face Hub')
    parser.add_argument('--split', type=str, default='train', help='Dataset split (default: train)')
    parser.add_argument('--output', type=str, help='Output JSONL file path')
    args = parser.parse_args()

    from datasets import load_dataset
    import json

    # Load the dataset
    dataset = load_dataset(args.dataset_name, split=args.split)
    
    # If output path not specified, use dataset name
    output_path = args.output or f"{args.dataset_name.replace('/', '_')}_{args.split}.jsonl"
    
    # Save as JSONL
    dataset.to_json(output_path)
    print(f"Dataset saved to {output_path}")

if __name__ == '__main__':
    main()
