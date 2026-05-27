"""Multi-platform Go binary release macro."""

load("@aspect_bazel_lib//lib:transitions.bzl", "platform_transition_filegroup")
load("@io_bazel_rules_go//go:def.bzl", "go_binary")

_PLATFORMS = [
    struct(suffix = "darwin_arm64", platform = "//platforms:darwin_arm64"),
    struct(suffix = "darwin_amd64", platform = "//platforms:darwin_amd64"),
    struct(suffix = "linux_amd64", platform = "//platforms:linux_amd64"),
    struct(suffix = "linux_arm64", platform = "//platforms:linux_arm64"),
    struct(suffix = "windows_amd64", platform = "//platforms:windows_amd64"),
]

def multi_platform_go_binaries(name, embed, visibility = []):
    """Creates a go_binary for each target platform.

    Args:
        name: name of the filegroup containing all platform-specific binaries.
        embed: go_library target(s) to embed into the binary.
        visibility: visibility for generated targets.
    """
    go_binary(
        name = "_{}".format(name),
        embed = embed,
        gc_linkopts = ["-s", "-w"],
        pure = "on",
        visibility = ["//visibility:private"],
    )

    targets = []
    for p in _PLATFORMS:
        target_name = "{}-{}".format(name, p.suffix)
        platform_transition_filegroup(
            name = target_name,
            srcs = [":_{}".format(name)],
            target_platform = p.platform,
            visibility = visibility,
        )
        targets.append(":{}".format(target_name))

    native.filegroup(
        name = name,
        srcs = targets,
        visibility = visibility,
    )
