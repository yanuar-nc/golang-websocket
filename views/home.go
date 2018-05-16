package views

const HomeHTML = `<!DOCTYPE html>
<html lang="en">
    <head>
        <title>WebSocket Example</title>
    </head>
    <body>
		<h1>
			<b>Total User:</b>
			<span  id="fileData">{{.Data}}</span>
		</h1>
        <script type="text/javascript">
            (function() {
                var data = document.getElementById("fileData");
                var conn = new WebSocket("ws://{{.Host}}/ws?lastMod={{.LastMod}}");
                conn.onclose = function(evt) {
                    data.textContent = 'Connection closed';
                }
                conn.onmessage = function(evt) {
                    console.log('file updated');
                    data.textContent = evt.data;
                }
            })();
        </script>
    </body>
</html>
`
