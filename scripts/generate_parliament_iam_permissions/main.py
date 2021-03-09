from generate_go_file import generate_go_file
from scrape_iam_permissions import scrape


def main():
    iam_permissions = scrape()
    generate_go_file(iam_permissions)
    print("Complete")


if __name__ == '__main__':
    main()
