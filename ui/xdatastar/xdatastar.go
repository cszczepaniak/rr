package xdatastar

import (
	"fmt"

	datastar "github.com/starfederation/datastar/sdk/go"
)

func ForceReloadScriptsOnPage(sse *datastar.ServerSentEventGenerator, elementQuery string) error {
	/* This is the minified version of the following script:

	function forceReloadScripts(node) {
		if (!node) {
			return
		}
		if (nodeScriptIs(node)) {
			node.parentNode.replaceChild(nodeScriptClone(node), node);
		} else {
			var i = -1, children = node.childNodes;
			while (++i < children.length) {
				  forceReloadScripts(children[i]);
			}
		}

		return node;
	}

	function nodeScriptClone(node){
		var script  = document.createElement("script");
		script.text = node.innerHTML;

		var i = -1, attrs = node.attributes, attr;
		while (++i < attrs.length) {
			  script.setAttribute((attr = attrs[i]).name, attr.value);
		}
		return script;
	}

	function nodeScriptIs(node) {
		if (node && node.tagName) {
			return node.tagName === 'SCRIPT';
		}
		return false;
	}

	forceReloadScripts(%s)

	*/
	script := fmt.Sprintf(`function forceReloadScripts(e){if(e){if(nodeScriptIs(e))e.parentNode.replaceChild(nodeScriptClone(e),e);else for(var t=-1,r=e.childNodes;++t<r.length;)forceReloadScripts(r[t]);return e}}function nodeScriptClone(e){var t=document.createElement("script");t.text=e.innerHTML;for(var r,n=-1,i=e.attributes;++n<i.length;)t.setAttribute((r=i[n]).name,r.value);return t}function nodeScriptIs(e){return!!e&&!!e.tagName&&"SCRIPT"===e.tagName}forceReloadScripts(%s);`, elementQuery)

	return sse.ExecuteScript(script)
}
