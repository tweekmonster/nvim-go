command! -nargs=* GoGuruCallees call GoGuru('callees', <f-args>)
command! -nargs=* GoGuruCallers call GoGuru('callers', <f-args>)
command! -nargs=* GoGuruCallstack call GoGuru('callstack', <f-args>)
command! -nargs=* GoGuruDefinition call GoGuru('definition', <f-args>)
command! -nargs=* GoGuruDescribe call GoGuru('describe', <f-args>)
command! -nargs=* GoGuruFreevars call GoGuru('freevars', <f-args>)
command! -nargs=* GoGuruImplements call GoGuru('implements', <f-args>)
command! -nargs=* GoGuruChannelPeers call GoGuru('peers', <f-args>)
command! -nargs=* GoGuruPointsto call GoGuru('pointsto', <f-args>)
command! -nargs=* GoGuruReferrers call GoGuru('referrers', <f-args>)
command! -nargs=* GoGuruWhicherrs call GoGuru('whicherrs', <f-args>)
