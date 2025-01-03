module.exports = {
  extends: ["@commitlint/config-conventional"],
  rules: {
    "scope-enum": [2, "always", ["python", "nodejs", "ruby", "global"]],
    "scope-empty": [2, "never"],
    "scope-case": [2, "always", "lower-case"],
    "subject-case": [2, "always", "lower-case"],
    "type-enum": [
      2,
      "always",
      [
        "feat",
        "fix",
        "docs",
        "style",
        "refactor",
        "perf",
        "test",
        "build",
        "ci",
        "chore",
        "revert",
      ],
    ],
  },
};
