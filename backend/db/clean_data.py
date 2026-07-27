#!/usr/bin/env python3
import csv
import sys
import unicodedata


def strip_accents(text):
    normalized = unicodedata.normalize("NFD", text)
    return "".join(c for c in normalized if not unicodedata.combining(c))


def clean_name(name):
    name = strip_accents(" ".join(name.split()))
    return " ".join(word[:1].upper() + word[1:] for word in name.split(" "))


def main(src_path, dst_path):
    seen_ids = set()
    with open(src_path, newline="", encoding="utf-8") as src, \
            open(dst_path, "w", newline="", encoding="utf-8") as dst:
        reader = csv.DictReader(src)
        writer = csv.writer(dst)
        writer.writerow(["id", "name"])

        for row in reader:
            original_id = row["id"].strip()
            name = clean_name(row["name"])

            if not original_id or not name or original_id in seen_ids:
                continue

            seen_ids.add(original_id)
            writer.writerow([original_id, name])


if __name__ == "__main__":
    main(sys.argv[1], sys.argv[2])
