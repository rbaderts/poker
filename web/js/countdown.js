'use strict';

var countdown =  Vue.component('countdown', {
	name: "countdown",
	props:{
		length:{
			type:Number,
			default:0
		},
	    startTime: {
			type:Number,
			default:0
		},
		elapsed: {
			type: Number,
			default: 0
		}
	},
	methods: {
	},
	computed: {
	},
	template: `
	<div class="timerbar"> 
	
	      <div v-bind:style="{remaining: activeColor}"></div>
	      
	</div>`
});

