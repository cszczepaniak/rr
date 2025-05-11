// This script allows us to force reloading JS after it gets swapped in by setting innerHTML.
// Otherwise, the script won't run when it's swapped in.
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
