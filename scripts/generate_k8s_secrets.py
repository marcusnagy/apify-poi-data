import base64
from pathlib import Path
from dotenv import dotenv_values
import argparse

def encode_to_base64(input_str):
    return base64.b64encode(input_str.encode()).decode()

def load_env_file(env_file_path):
    if not Path(env_file_path).exists():
        raise FileNotFoundError(f"The file {env_file_path} does not exist.")

    return dotenv_values(env_file_path)

def generate_k8s_secret_format(env_vars):
    lines = []
    for key, value in env_vars.items():
        encoded_value = encode_to_base64(value)
        lines.append(f"{key}: {encoded_value}")
    return "\n".join(lines)

def generate_env_list_format(env_vars):
    keys = [key for key in env_vars.keys()]
    return "\n".join(f"- {key}" for key in keys)

def main():
    parser = argparse.ArgumentParser(description="Generate Kubernetes secrets and environment variable list from a .env file.")
    parser.add_argument("env_file", type=str, help="Path to the .env file")
    args = parser.parse_args()

    try:
        env_vars = load_env_file(args.env_file)
        print("\nKubernetes Secret Format:")
        print(generate_k8s_secret_format(env_vars))

        print("\nEnvironment Variable List:")
        print(generate_env_list_format(env_vars))
    except FileNotFoundError as e:
        print(e)
    except Exception as e:
        print(f"An unexpected error occurred: {e}")

if __name__ == "__main__":
    main()
