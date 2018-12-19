// SPDX-License-Identifier: Apache-2.0

'use strict';

var app = angular.module('application', []);

// Angular Controller
app.controller('appController', function($scope, appFactory){

	$("#success_holder").hide();
	$("#success_create").hide();
	$("#error_holder").hide();
	$("#error_query").hide();

	$scope.queryResume = function(){

		var id = $scope.resume_id;

		appFactory.queryResume(id, function(data){
			$scope.query_resume = data;

			if ($scope.query_resume == "Could not locate resume"){
				console.log()
				$("#error_query").show();
			} else{
				$("#error_query").hide();
			}
		});
	}

	$scope.recordResume = function(){

		appFactory.recordResume($scope.resume, function(data){
			$scope.create_resume = data;
			$("#success_create").show();
		});
	}

});

// Angular Factory
app.factory('appFactory', function($http){
	
	var factory = {};

	factory.queryResume = function(id, callback){
    	$http.get('/get_resume/'+id).success(function(output){
			callback(output)
		});
	}

	factory.recordResume = function(data, callback){

		var resume = data.resume + "-" + data.user_id;

    	$http.get('/add_resume/'+resume).success(function(output){
			callback(output)
		});
	}

	return factory;
});


