# 15 december 2019
dir=$(mktemp -d) || exit $?
[ "Z$savedir" = Z1 ] && echo "$dir"
GOCACHE="$dir" go "$@"
status=$?
[ "Z$savedir" = Z1 ] || rm -rf "$dir"
exit $status
