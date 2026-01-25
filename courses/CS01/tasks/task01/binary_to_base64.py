# Tâche: convertisseur d'un binaire quelconque en Base64.
# Description:  implémentez une fonction encodeBase64(N) qui prend un nombre binaire et qui retourne un string reprenant sa représentation Base64 (attention au Padding)
import argparse
import base64


def binary_to_base64(binary_str: str) -> str:
    """Convert a binary string to its Base64 representation."""
    # Pad the binary string to make its length a multiple of 8
    while len(binary_str) % 8 != 0:
        binary_str = "0" + binary_str

    # Convert binary string to bytes
    byte_array = bytearray()
    for i in range(0, len(binary_str), 8):
        byte = binary_str[i : i + 8]
        byte_array.append(int(byte, 2))

    # Encode bytes to Base64
    base64_bytes = base64.b64encode(byte_array)
    return base64_bytes.decode("utf-8")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Convert a binary string to its Base64 representation."
    )
    parser.add_argument("number", type=str, help="The binary string to convert.")
    args = parser.parse_args()
    print(binary_to_base64(args.number))
