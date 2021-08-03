'use strict';

Vue.component('table-icon', {
	// The todo-item component now accepts a
	// "prop", which is like a custom attribute.
	// This prop is called todo.
	props: ['table'],
	template: `
	      <div class="table-icon">
	      <li>
<!--               <button style="background-image: /static/img/poker_table_icon.png" class="table_button" v-on:click="join()"/>-->
               <input type="image" src="/static/img/poker_table_icon.png" class="table_button btTxt submit" v-on:click="join()"/>

	      </li>
	      </div>`,
	methods:
		{
			join: function () {
//				var id = this.$vnode.data.key;
				var thisvm = this;
				var id = thisvm.table.id;
				//console.log("key = " + this.$vnode.data.key);
				console.log("key = " + id);
				console.log("name " + thisvm.table.name + " clicked");
				console.log("id " + thisvm.table.id + " clicked");

				//var tableId = window.location.pathname.split("/")[2];

				var joinurl = "/api/tables/" + id + "/join";
				var url = "tables/" + id;
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
			//			let tableData = data["table"];
						let jwt = data["jwt"];
						window.jwt = jwt;
						let newurl = url;
						let newwindow = window.open(newurl, "_blank");
						newwindow.self.jwt = jwt;
						newwindow.self.tableId = id;
//						newwindow.self.table = tableData;

						//newwindow.location.reload();
						//gBrowser.addTab("http://www.google.com/");
						//gBrowser.selectedTab = gBrowser.addTab("http://www.google.com/");


						//window.location.assign(url);
						//tableResult = data;
						console.log("data = " + JSON.stringify(data))
					}
				});

			}
		}



});
