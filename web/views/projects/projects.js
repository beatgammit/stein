define([
	'angular',
], function (angular) {
	'use strict';

	angular.module('stein.projects', [])
		.controller('ProjectsCtrl', ['$scope', '$http', function ($scope, $http) {
			$http.get('/projects').success(function (data) {
				$scope.projects = data;
			});
		}]);
});
