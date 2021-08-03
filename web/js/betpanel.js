
'use strict';


var betpanel =  Vue.component('betpanel', {
	name: "betpanel",
	props: {
		requestId: {
			type: Number,
			default: 0
		},
		currentBet: {
			type: Number,
			default: 0
		},
		callAmount: {
			type: Number,
			default: 0
		},
		stackAvailable: {
			type: Number,
			default: 0
		}
	},


	/*
<div className="card" style="width: 18rem;">

	<img className="card-img-top" src="..." alt="Card image cap">
		<div className="card-body">
			<h5 className="card-title">Card title</h5>
			<p className="card-text">Some quick example text to build on the card title and make up the bulk of the
				card's content.</p>
			<a href="#" className="btn btn-primary">Go somewhere</a>
		</div>
</div>

*/
template: `
     	<div class="betinput">
     	
          	<b-container class="betinput-container mr-2">
          	<!--
                <b-form-row class="my-1">
                    <b-col sm="12"><h5>Amount to call: {{callAmount}}</h5></b-col>
			    </b-form-row>
			    -->
			    
                <b-row class="my-1 no-gutters">
                    <b-col md="2" class="mx-1 px-0 col">
						<b-button variant="danger" class="" size="sm" v-on:click="fold" type="button">Fold</b-button>
                    </b-col>
                    <b-col v-if="callAmount == 0" md="2" class="ml-1 px-0 col">
     					<b-button variant="primary" class="" size="sm" v-on:click="check" type="button">Check</b-button>
                    </b-col>
                    <b-col md="2" class="mx-1 px-0 col" offset-md="2">
						<b-button variant="warning" class="" size="sm" v-on:click="allIn" type="button">All In</b-button>
                    </b-col>
				</b-row>
                <b-row class="my-2 no-gutters">
                    <b-col md="3" class="col mx-1 px-0">
     					<b-button variant="danger" class="" size="sm" v-on:click="callsubmit" type="button">Call
       					    <b-badge class="ml-0 mt-0" variant="light"> {{ callAmount }} <span class="sr-only">call amount</span></b-badge>
       					    </b-button>
                    </b-col>
                    <b-col md="1" class="col mx-1 px-0">
      					<b-form-input id="amount-1" size="sm" v-model="bet" type="number" min="0"
      					 :max="stackAvailable"
      					 style="width=35px; padding:2px;"></b-form-input> 
                    </b-col>
                    <b-col md="4" class="col mx-1 mr-0 px-0 pt-1">
	     				<b-form-input id="amount-slider"  size="sm" v-model="bet" type="range" min="0" :max="stackAvailable"></b-form-input>
	     				<!--<b-form-input id="amount-slider" v-on:change="slide($event)"  size="sm" v-model="bet" type="range" min="0" :max="stackAvailable"></b-form-input>-->
                    </b-col>
                    <b-col md="2" class="col mx-1 px-0">
     		            <b-button variant="danger" id="bet-button" size="sm" class="" v-on:click="betsubmit" type="button">Bet X</b-button>
                    </b-col>
				</b-row>
				<!--
                <b-row class="my-1">
                    <b-col md="3" class="mx-0 px-0">
     		            <b-button id="bet-button" size="sm" class="" v-on:click="betsubmit" type="button">Bet X</b-button>
                    </b-col>
                    <b-col md="3" offset-md="1" class="mx-0 px-0">
      					<b-form-input id="amount-1" size="sm" v-model="bet" type="number" min="0" :max="stackAvailable"></b-form-input> 
                    </b-col>
                    <b-col md="5" class="ml-1 mr-0 px-0">
	     				<b-form-input id="amount-slider"  size="sm" v-model="bet" type="range" min="0" :max="stackAvailable"></b-form-input>
                    </b-col>
   			    </b-row>
   			    -->
          </b-container>

	<!--
							<button v-on:click="betsubmit" id="bet-submit-button" type="button" 
								   class="btn btn-primary">Bet</button>
								   -->

<!--
						  <input class="form-control form-control-sm" id="bet-amount" 
								 style="width: 40px"
								 v-bind:bet="bet" ref="input" 
								 v-on:input="$emit('input', $event.target.bet)"
								 aria-describedby="bet-submit-button">
							<button v-on:click="betsubmit" id="bet-submit-button" type="button" 
								   class="btn btn-primary">Bet</button>
						   </input>
						   -->
				   
				   <!--<b-progress :value="value" :max="max" class="mb-3"></b-progress> -->


<!--
				   <div class="form-group form-group-sm row">
					   <button v-on:click="betsubmit" id="bet-submit-button" type="button" class="btn btn-primary">Post Bet</button>
					   <button v-on:click="callsubmit" type="button" class="btn btn-primary">Call</button>
					   <button v-on:click="raise2submit" type="button" class="btn btn-primary">Raise X 2</button>
				   </div>
				   -->

	   </div>`,

	data() {
		return {
			startTime: 0,
			timeremaining: 500000,
			interval: null,
			bet: 0
		}
	},
	computed: {
		// a computed getter
		enteredbet: function () {
			// `this` points to the vm instance
			return this.bet;
		}
	},
	mounted: function() {
	    let now = new Date();
     	this.startTime = now.getTime();
     	this.startTimer();
	},
	methods : {

	    startTimer: function() {
	    	var vm = this;

		    this.interval = setInterval(function() {
		    	let now = new Date();
		    	let elapsed = now.getTime() - vm.startTime;
		    	vm.timeremaining = 500000-elapsed;

		    	if (elapsed > 500000) {
		    		vm.fold(null);
				    vm.$parent.showBetPanel = false;
		    		clearInterval(vm.interval)
			    }
		    }, 1000);

	    },

		allIn: function(event) {
			var vm = this;
			vm.sendBet(event,vm.stackAvailable, 1, 0);
		},

		fold: function(event) {
			var vm = this;
			vm.sendBet(event, 0, 0, 1);
		},
        check: function(event) {
	        var vm = this;
			vm.sendBet(event, 0, 0, 0);
        },
	    raise2submit: function(event) {
		    var vm = this;
			console.log("raise X 2");
			vm.sendBet(event, (vm.callAmount * 2), 0, 0);
	    },
		callsubmit: function(event) {
			var vm = this;
			console.log("call");
			vm.sendBet(event, vm.callAmount, 0, 0);
		},
		betsubmit: function(event) {
			var vm = this;
			let newbet = vm.enteredbet;
			vm.sendBet(event, newbet, 0, 0);
		},
		sendBet: function(event, amount, isAllIn, isFold) {
			var vm = this;
			console.log("bet = " + amount);
			ActionCounter += 1;
			let action = {
				Id: ActionCounter,
				ResponseTo: Number(vm.requestId),
				data: {"amount": Number(amount),
					   "isAllIn" : isAllIn,
     				   "isFold" : isFold}
			};

			console.log(JSON.stringify(action));
			let id = window.tableId;
			axios.post(`/api/tables/`+id+`/bet`, action, {withCredentials: false} );

			vm.$parent.showBetPanel = false;
			clearInterval(vm.interval);
			vm.callAmount = 0;
			vm.currentBet = 0;

			//vm.$emit('bet', Number(amount));
		},
		slide: function(value) {
	    	vm.bet = value;
			let betbutton = document.getElementById('bet-button');
			betbutton.innerHTML = Bet + value;
		}
	}
});

