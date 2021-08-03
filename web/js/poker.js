'use strict';



/*
function joinTable(tableId) {

	var table = new Table(tableId);
	var WS = new WebSocket("ws://" + window.location.host + "/table/" + tableId);
	WS.binaryType = "arraybuffer";

	WS.onopen = function () {
		console.log("Server connection enstablished");
	};

	WS.onclose = function () {
		console.log("Server close");
		WS = null;
		var page = "/lobby";
		window.location.assign(page);
	};

	WS.onerror = function (error) {
		console.log("Server error: " + error);
		WS = null;
		var page = "/lobby";
		window.location.assign(page);
	};

	WS.onmessage = function (msg) {
		handleTableCommand(msg);
	};
}

function handleTableCommand(msg) {

}

 */



Vue.component('table-icon', {
	// The todo-item component now accepts a
	// "prop", which is like a custom attribute.
	// This prop is called todo.
	props: ['table'],
	template: `
	      <div class="table-icon">
	      <li>
               <button class="table_button" v-on:click="join(key)"/>
<!--                     <button v-on:click="$emit(\\'remove\\')">Remove</button>\\-->
	      </li>
	      </div>`,
	methods:
		{
			join: function (table) {
				var id = this.$vnode.data.key;
				console.log("key = " + this.$vnode.data.key);
//				console.log("key = " + table.Id);
//				console.log("key = " + $key);
				console.log("name " + this.name + " clicked");
				console.log("id " + this.id + " clicked");

				var joinurl = "/tables/" + id + "/join";
				var url = "/tables/" + id;
				var tableResult = null;
				$.ajax
				({
					type: "GET",
					//the url where you want to sent the userName and password to
					url: joinurl,
					dataType: 'json',
					async: false,
					data: null,
					success: function (data) {
						window.location.assign(url);
						tableResult = data;
						console.log("data = " + JSON.stringify(data))
					}
				});


//				var table = new Table(this.$vonde.data.key);
				/*
				var WS = new WebSocket("ws://" + window.location.host + "/tables/" + id + "/updates");
				WS.binaryType = "arraybuffer";

				WS.onopen = function () {
					console.log("Server connection enstablished");
				};

				WS.onclose = function () {
					console.log("Server close");
					WS = null;
					var page = "/lobby";
					window.location.assign(page);
				};

				WS.onerror = function (error) {
					console.log("Server error: " + error);
					WS = null;
					var page = "/lobby";
					window.location.assign(page);
				};

				WS.onmessage = function (msg) {
					handleTableCommand(msg);

				};

				 */
			}
		}
});


