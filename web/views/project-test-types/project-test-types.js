define([
	'angular',
], function (angular) {
	'use strict';

	angular.module('stein.project-test-types', [])
		.controller('ProjectTestTypesCtrl', ['$scope', '$routeParams', '$http', function ($scope, $routeParams, $http) {
			$scope.project = $routeParams.project;
			$http.get('/projects/' + $routeParams.project + '/types').success(function (data) {
				$scope.testTypes = data;
			});
		}]);
});
