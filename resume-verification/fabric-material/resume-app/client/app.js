// SPDX-License-Identifier: Apache-2.0

'use strict';

var app = angular.module('application', []);

// Angular Controller
app.controller('appController', function($scope, appFactory){

	$("#success_holder").hide();
	$("#success_create").hide();
	$("#error_holder").hide();
	$("#error_query").hide();
	
	$scope.queryAllResumes = function(){

		appFactory.queryAllResumes(function(data){
			var array = [];
			for (var i = 0; i < data.length; i++){
				parseInt(data[i].Key);
				data[i].Record.Key = parseInt(data[i].Key);
				array.push(data[i].Record);
			}
			array.sort(function(a, b) {
			    return parseFloat(a.Key) - parseFloat(b.Key);
			});
			$scope.all_resumes = array;
		});
	}

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

    factory.queryAllResumes = function(callback){

    	$http.get('/get_all_resumes/').success(function(output){
			callback(output)
		});
	}

	factory.queryResume = function(id, callback){
    	$http.get('/get_resume/'+id).success(function(output){
			callback(output)
		});
	}

	factory.recordResume = function(data, callback){

		var resume = data.id + "-" + data.resume + "-" + data.timestamp + "-" + data.resume_hash + "-" + data.user_id;

    	$http.get('/add_resume/'+resume).success(function(output){
			callback(output)
		});
	}

	return factory;
});


