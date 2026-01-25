import argparse


def decimal_to_binary(n: int) -> str:
    """Convert a decimal integer to its binary representation as a string."""
    if n == 0:
        return "0"

    binary_digits = []
    while n > 0:
        remainder = n % 2
        binary_digits.append(str(remainder))
        n = n // 2

    # The digits are in reverse order, so we need to reverse them before joining
    binary_digits.reverse()
    return "".join(binary_digits)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Convert a decimal integer to its binary representation."
    )
    parser.add_argument("number", type=int, help="The decimal integer to convert.")
    args = parser.parse_args()
    print(decimal_to_binary(args.number).zfill(16))
