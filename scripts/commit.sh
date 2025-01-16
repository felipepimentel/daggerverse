#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Valid commit types
TYPES=("feat" "fix" "perf" "docs" "style" "refactor" "test" "build" "ci" "chore")

# Language-specific scopes
LANGUAGE_SCOPES=("python" "nodejs" "ruby")

# Help message
show_help() {
    echo -e "${GREEN}Commit Message Format Helper${NC}"
    echo
    echo "Usage: $0 [options]"
    echo
    echo "This script helps format commit messages according to the repository's conventions."
    echo
    echo "Valid commit types:"
    for type in "${TYPES[@]}"; do
        echo "  - $type"
    done
    echo
    echo "Language-specific scopes:"
    for scope in "${LANGUAGE_SCOPES[@]}"; do
        echo "  - $scope"
    done
    echo
    echo "Options:"
    echo "  -h, --help     Show this help message"
    echo "  -t, --type     Specify commit type"
    echo "  -s, --scope    Specify commit scope"
    echo "  -m, --message  Specify commit message"
    echo "  -b, --body     Specify commit body"
    echo "  -f, --footer   Specify commit footer"
    echo
    echo "Example:"
    echo "  $0 -t feat -s python-pipeline -m 'add new testing feature' -b 'Implements comprehensive testing' -f 'BREAKING CHANGE: New API structure'"
}

# Function to validate commit type
validate_type() {
    local type=$1
    for valid_type in "${TYPES[@]}"; do
        if [ "$type" == "$valid_type" ]; then
            return 0
        fi
    done
    return 1
}

# Function to validate scope
validate_scope() {
    local scope=$1
    
    # Check if it's a language-specific scope
    for lang in "${LANGUAGE_SCOPES[@]}"; do
        if [[ "$scope" == "$lang-"* ]]; then
            return 0
        fi
    done
    
    # Check if it's "global"
    if [ "$scope" == "global" ]; then
        return 0
    fi
    
    return 1
}

# Function to validate subject line length
validate_subject_length() {
    local subject=$1
    if [ ${#subject} -gt 72 ]; then
        return 1
    fi
    return 0
}

# Initialize variables
TYPE=""
SCOPE=""
MESSAGE=""
BODY=""
FOOTER=""
BREAKING=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -t|--type)
            TYPE="$2"
            shift 2
            ;;
        -s|--scope)
            SCOPE="$2"
            shift 2
            ;;
        -m|--message)
            MESSAGE="$2"
            shift 2
            ;;
        -b|--body)
            BODY="$2"
            shift 2
            ;;
        -f|--footer)
            FOOTER="$2"
            shift 2
            ;;
        *)
            echo -e "${RED}Error: Unknown option $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

# Interactive mode if no arguments provided
if [ -z "$TYPE" ] || [ -z "$SCOPE" ] || [ -z "$MESSAGE" ]; then
    echo -e "${GREEN}Interactive Commit Message Creation${NC}"
    echo

    # Get commit type
    if [ -z "$TYPE" ]; then
        echo -e "${YELLOW}Available commit types:${NC}"
        printf "%s " "${TYPES[@]}"
        echo
        read -p "Enter commit type: " TYPE
    fi

    # Get scope
    if [ -z "$SCOPE" ]; then
        echo -e "\n${YELLOW}Enter scope (e.g., python-pipeline, nodejs-app, global):${NC}"
        read -p "Enter scope: " SCOPE
    fi

    # Get commit message
    if [ -z "$MESSAGE" ]; then
        echo -e "\n${YELLOW}Enter commit message:${NC}"
        read -p "Enter message: " MESSAGE
    fi

    # Get commit body (optional)
    if [ -z "$BODY" ]; then
        echo -e "\n${YELLOW}Enter commit body (optional, press Enter to skip):${NC}"
        read -p "Enter body: " BODY
    fi

    # Get commit footer (optional)
    if [ -z "$FOOTER" ]; then
        echo -e "\n${YELLOW}Enter commit footer (optional, press Enter to skip):${NC}"
        read -p "Enter footer: " FOOTER
    fi
fi

# Validate commit type
if ! validate_type "$TYPE"; then
    echo -e "${RED}Error: Invalid commit type '$TYPE'${NC}"
    echo "Valid types are: ${TYPES[*]}"
    exit 1
fi

# Validate scope
if ! validate_scope "$SCOPE"; then
    echo -e "${RED}Error: Invalid scope '$SCOPE'${NC}"
    echo "Scope must be 'global' or start with one of: ${LANGUAGE_SCOPES[*]}"
    exit 1
fi

# Validate subject line length
if ! validate_subject_length "$MESSAGE"; then
    echo -e "${RED}Error: Commit message is too long (max 72 characters)${NC}"
    exit 1
fi

# Check if message starts with uppercase letter
if [[ ! "$MESSAGE" =~ ^[A-Z] ]]; then
    echo -e "${RED}Error: Commit message must start with an uppercase letter${NC}"
    exit 1
fi

# Check if message ends with a period
if [[ "$MESSAGE" =~ \.$  ]]; then
    echo -e "${RED}Error: Commit message should not end with a period${NC}"
    exit 1
fi

# Check for breaking changes in footer
if [[ "$FOOTER" == *"BREAKING CHANGE:"* ]]; then
    BREAKING=true
fi

# Construct the commit message
COMMIT_MSG="$TYPE($SCOPE)"
if [ "$BREAKING" = true ]; then
    COMMIT_MSG="$COMMIT_MSG!"
fi
COMMIT_MSG="$COMMIT_MSG: $MESSAGE"

if [ ! -z "$BODY" ]; then
    COMMIT_MSG="$COMMIT_MSG

$BODY"
fi

if [ ! -z "$FOOTER" ]; then
    COMMIT_MSG="$COMMIT_MSG

$FOOTER"
fi

# Show the final commit message
echo -e "\n${GREEN}Final commit message:${NC}"
echo -e "${YELLOW}-------------------${NC}"
echo "$COMMIT_MSG"
echo -e "${YELLOW}-------------------${NC}"

# Confirm and commit
read -p "Do you want to commit with this message? (y/N) " confirm
if [[ $confirm =~ ^[Yy]$ ]]; then
    # Check if there are staged changes
    if git diff --cached --quiet; then
        echo -e "${RED}Error: No changes staged for commit${NC}"
        echo "Use 'git add' to stage changes before committing"
        exit 1
    fi
    
    # Perform the commit
    echo "$COMMIT_MSG" | git commit -F -
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}Successfully committed changes!${NC}"
    else
        echo -e "${RED}Failed to commit changes${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}Commit cancelled${NC}"
    exit 0
fi 