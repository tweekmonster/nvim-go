if exists('g:loaded_nvim_go')
  finish
endif
let g:loaded_nvim_go = 1

let s:plugin_name = 'nvim-go'
let s:goos = split(system('go env GOOS'), "\n")[0]
let s:goarch = split(system('go env GOARCH'), "\n")[0]
let s:plugin_path = fnamemodify(resolve(expand('<sfile>:p')), ':h:h') . '/bin/' . s:plugin_name . '-' . s:goos . '-' . s:goarch

function! s:RequireGoHost(host) abort
  let args = []
  try
    for plugin in remote#host#PluginsForHost(a:host.name)
        call add(args, plugin.path)
    endfor
    return rpcstart(s:plugin_path, args)
  catch
    echomsg v:exception
  endtry
  throw 'Failed to load ' . s:plugin_name . ' host'.
endfunction

call remote#host#Register('nvim-go', '*', function('s:RequireGoHost'))