/* global describe, beforeEach, inject, it, expect */
'use strict';

describe('controllers', function(){
  var scope;

  beforeEach(module('js'));

  beforeEach(inject(function($rootScope) {
  	scope = $rootScope.$new();
  }));

  it('Should define two functions', inject(function($controller) {
    expect(scope.generate).toBeUndefined();
    expect(scope.send).toBeUndefined();

    $controller('MainCtrl', {
      $scope: scope
  	});

    expect(angular.isFunction(scope.generate)).toBeTruthy();
    expect(angular.isFunction(scope.send)).toBeTruthy();
  }));
});
