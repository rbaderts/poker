'use strict';


var card =  Vue.component('playerbet', {
	name: "playerbet",
	template: `<div class="mise">
    			    <div class="mise-value">{{ amount }}</div>
					<div class="jeton-10">
						<div class="jetons v-10" v-for="(n, i) in ((player.onTable - (player.onTable % 10)) / 10)" :style="{top: (-2 + (i * 5)) + 'px'}" v-if="player.onTable / 10 >= 1"></div>
					</div>
					<div class="jeton-5">
						<div class="jetons v-5" v-for="(n, i) in (((player.onTable % 10) - ((player.onTable % 10) % 2)) / 2)" :style="{top: (-2 + (i * 5)) + 'px'}" v-if="player.onTable % 10 && player.onTable % 10 >= 2"></div>
					</div>
					<div class="jeton-1">
						<div class="jetons v-1" v-if="player.onTable % 10 && player.onTable % 2"></div>
					</div>
				</div>`,
	props: {
		amount: {
			type: Number,
			default: 0
		},
	},
	computed: {
	} ,
	methods : {
	}
});
/*
var card =  Vue.component('card', {
	name: "card",
	template: `
		<div class="card" :class="[suitSym, 'values-' + indexSym]"
		style="background-image: "image">
			<h1>{{indexSym}}</h1>
			<div class="figures" :class="suitSym"></div>
			<h1>{{indexSym}}</h1>
		</div> `,
	props: ['index', 'suit'],
	computed: {
		indexSym: function () {
			return this.cardSymbol(this.index);
		},
		suitSym: function () {
			return this.cardSuit(this.suit);
		},

	},
	methods: {
		cardSuit: function(suit) {
			if (suit == 1) {
				return 'figures-P';
			} else if (suit == 2) {
				return 'figures-H';
			} else if (suit == 3) {
				return 'figures-D';
			} else if (suit == 4) {
				return 'figures-C';
			}
		},
		cardSymbol: function(index) {
			if (index < 10) {
				return index
			} else if (index == 11) {
				return 'J';
			} else if (index == 12) {
				return 'J';
			} else if (index == 13) {
				return 'K';
			} else if (index == 14) {
				return 'A';
			}
		}


	}
});
*/
