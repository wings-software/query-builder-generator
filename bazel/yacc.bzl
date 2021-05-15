load("@io_bazel_rules_go//go:def.bzl", "GoSource",)

_GO_YACC_TOOL = "@org_golang_x_tools//cmd/goyacc"

def _go_yacc_impl(ctx):
    args = ctx.actions.args()
    args.add("-o", ctx.outputs.out)
    args.add("-p", ctx.attr.prefix)
    args.add(ctx.file.src)
    goroot = "%s/.." % ctx.executable._go_yacc_tool.dirname
    ctx.actions.run(
        executable = ctx.executable._go_yacc_tool,
        arguments = [args],
        inputs = [ctx.file.src],
        outputs = [ctx.outputs.out],
        env = {
            "GOROOT": goroot,
        },
    )
    return DefaultInfo(
        files = depset([ctx.outputs.out]),
    )

_go_yacc = rule(
    implementation = _go_yacc_impl,
    attrs = {
        "src": attr.label(
            allow_single_file = True,
        ),
        "out": attr.output(),
        "prefix": attr.string(
            default='yy',
            doc='',
            mandatory=False
        ),
        "_go_yacc_tool": attr.label(
            default = _GO_YACC_TOOL,
            allow_single_file = True,
            executable = True,
            cfg = "host",
        ),
    },
)

def go_yacc(yacc, go, prefix, visibility = None):
    _go_yacc(
        name = yacc + ".go_yacc",
        src = yacc,
        out = go,
        prefix = prefix,
        visibility = visibility,
    )