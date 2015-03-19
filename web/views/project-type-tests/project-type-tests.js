define([
	'angular',
], function (angular) {
	'use strict';

	angular.module('stein.project-type-tests', [])
		.controller('ProjectTypeTestsCtrl', ['$scope', '$routeParams', '$http', function ($scope, $routeParams, $http) {
			$scope.project = $routeParams.project;
			$scope.type = $routeParams.type;
			$http.get('/projects/' + $routeParams.project + '/types/' + $routeParams.type).success(function (data) {
				$scope.tests = data;
			});
		}]);
});
