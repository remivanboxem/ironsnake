# Tâche: convertisseur d'un binaire quelconque en Base64.
# Description:  implémentez une fonction encodeBase64(N) qui prend un nombre binaire et qui retourne un string reprenant sa représentation Base64 (attention au Padding)
import argparse


@@binary_to_base64@@

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Convert a binary string to its Base64 representation."
    )
    parser.add_argument("number", type=str, help="The binary string to convert.")
    args = parser.parse_args()
    print(binary_to_base64(args.number))
