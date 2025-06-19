#!/bin/bash

# Version management script for replyzer

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Function to show current version
show_version() {
    local current_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "No tags found")
    local commit_count=$(git rev-list --count HEAD)
    local commit_hash=$(git rev-parse --short HEAD)
    
    echo "Current version: $current_tag"
    echo "Total commits: $commit_count"
    echo "Current commit: $commit_hash"
    
    if [ "$current_tag" != "No tags found" ]; then
        echo "Commits since last tag: $(git rev-list --count ${current_tag}..HEAD)"
    fi
}

# Function to show what the next version would be
preview_next_version() {
    echo "Analyzing commits since last release..."
    
    local last_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
    local commits_range
    
    if [ -n "$last_tag" ]; then
        commits_range="${last_tag}..HEAD"
    else
        commits_range="HEAD"
    fi
    
    local has_feat=$(git log --oneline $commits_range | grep -E "^[a-f0-9]+ feat" | wc -l)
    local has_fix=$(git log --oneline $commits_range | grep -E "^[a-f0-9]+ fix" | wc -l)
    local has_breaking=$(git log --oneline $commits_range | grep -E "BREAKING CHANGE|!" | wc -l)
    
    echo "Changes since last release:"
    echo "  Features (minor): $has_feat"
    echo "  Bug fixes (patch): $has_fix"
    echo "  Breaking changes (major): $has_breaking"
    
    if [ $has_breaking -gt 0 ]; then
        echo "Next version would be: MAJOR (breaking changes detected)"
    elif [ $has_feat -gt 0 ]; then
        echo "Next version would be: MINOR (new features detected)"
    elif [ $has_fix -gt 0 ]; then
        echo "Next version would be: PATCH (bug fixes detected)"
    else
        echo "Next version would be: No release (no significant changes)"
    fi
}

# Function to validate commit messages
validate_commits() {
    echo "Validating commit messages since last release..."
    
    local last_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
    local commits_range
    
    if [ -n "$last_tag" ]; then
        commits_range="${last_tag}..HEAD"
    else
        commits_range="HEAD"
    fi
    
    local invalid_commits=0
    
    while IFS= read -r commit; do
        local commit_hash=$(echo "$commit" | cut -d' ' -f1)
        local commit_msg=$(echo "$commit" | cut -d' ' -f2-)
        
        if ! echo "$commit_msg" | grep -qE "^(feat|fix|docs|style|refactor|test|chore|perf|ci|build|revert)(\(.+\))?\!?:"; then
            echo "❌ Invalid commit: $commit_hash - $commit_msg"
            invalid_commits=$((invalid_commits + 1))
        else
            echo "✅ Valid commit: $commit_hash - $commit_msg"
        fi
    done < <(git log --oneline $commits_range)
    
    if [ $invalid_commits -gt 0 ]; then
        echo ""
        echo "⚠️  Found $invalid_commits invalid commit message(s)."
        echo "Please use conventional commit format: type(scope): description"
        return 1
    else
        echo ""
        echo "✅ All commit messages are valid!"
        return 0
    fi
}

# Main script logic
case "$1" in
    "show"|"")
        show_version
        ;;
    "preview")
        preview_next_version
        ;;
    "validate")
        validate_commits
        ;;
    "help")
        echo "Usage: $0 [show|preview|validate|help]"
        echo ""
        echo "Commands:"
        echo "  show (default) - Show current version information"
        echo "  preview        - Preview what the next version would be"
        echo "  validate       - Validate commit messages for semantic versioning"
        echo "  help          - Show this help message"
        ;;
    *)
        echo "Unknown command: $1"
        echo "Use '$0 help' for usage information"
        exit 1
        ;;
esac