define([
	'angular',
], function (angular) {
	'use strict';

	angular.module('stein.project-tests', [])
		.controller('ProjectTestsCtrl', ['$scope', '$routeParams', '$http', function ($scope, $routeParams, $http) {
			$scope.project = $routeParams.project;
			$http.get('/projects/' + $routeParams.project + '/tests').success(function (data) {
				$scope.tests = data;
			});
		}]);
});
