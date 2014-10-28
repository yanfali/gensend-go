/* global describe, beforeEach, inject, it, expect */
'use strict';

describe('controllers-generate', function() {
  var scope;

  beforeEach(module('js'));

  beforeEach(inject(function($rootScope) {
    scope = $rootScope.$new();
  }));

  it('Should define a json form', inject(function($controller) {
    expect(scope.generateForm).toBeUndefined();

    $controller('GenerateCtrl', {
      $scope: scope
    });

    expect(angular.isDefined(scope.generateForm)).toBeTruthy();
  }));

  it('json form Must have members', inject(function($controller) {
    expect(scope.generateForm).toBeUndefined();

    $controller('GenerateCtrl', {
      $scope: scope
    });

    (function(scope) {
      var form = scope.generateForm;
      ['length', 'letters', 'mixedCase', 'numbers', 'punctuation', 'similarChars'].forEach(function(value) {
        expect(angular.isDefined(form[value])).toBeTruthy('missing ' + value);
      });
    }(scope));
  }));
});
