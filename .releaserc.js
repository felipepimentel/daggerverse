module.exports = {
  branches: ["main"],
  tagFormat: "${MODULE_NAME}/v${version}",
  plugins: [
    [
      "@semantic-release/commit-analyzer",
      {
        preset: "angular",
        releaseRules: [
          { type: "feat", scope: "python", release: "minor" },
          { type: "fix", scope: "python", release: "patch" },
          { type: "perf", scope: "python", release: "patch" },
          { type: "docs", scope: "python", release: "patch" },
          { breaking: true, release: "major" },
        ],
        parserOpts: {
          noteKeywords: ["BREAKING CHANGE", "BREAKING CHANGES"],
        },
      },
    ],
    "@semantic-release/release-notes-generator",
    [
      "@semantic-release/changelog",
      {
        changelogFile: "${MODULE_PATH}/CHANGELOG.md",
      },
    ],
    [
      "@semantic-release/git",
      {
        assets: ["${MODULE_PATH}/CHANGELOG.md"],
        message:
          "chore(${MODULE_NAME}): release ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}",
      },
    ],
    [
      "@semantic-release/github",
      {
        successComment:
          "🎉 This PR is included in version ${nextRelease.version}",
        failTitle: "The release workflow failed",
        failComment:
          "The release workflow failed. Please check the logs for more details.",
        releasedLabels: ["released"],
        addReleases: "bottom",
      },
    ],
  ],
};
