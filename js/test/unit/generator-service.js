/* global describe, beforeEach, inject, it, expect */
'use strict';

describe('service', function() {
  var generatorService;

  beforeEach(module('generator'));
  beforeEach(inject(function(_generatorService_) {
    generatorService = _generatorService_;
  }));

  it('Should offer getRandomUint32 service', function() {
    var array = generatorService.getRandomUint32(4);
    expect(angular.isDefined(array)).toBeTruthy();
  });

  it('getRandomUint32 should return 4 32bit words', function() {
    var array = generatorService.getRandomUint32(4);
    expect(array.length).toEqual(4);
  });

  it('getRandomUint32 should return 32bit words', function() {
    var array = generatorService.getRandomUint32(4);
    for (var i in [0, 1, 2, 3]) {
      expect(array[i]).toBeGreaterThan(0);
      expect(array[i]).toBeLessThan(Math.pow(2, 32) - 1);
    }
  });

  it('full lowercase set', function() {
    var testStr = generatorService.FULL_LOWERCASE;
    expect(angular.isDefined(testStr)).toBeTruthy();
    expect(angular.isString(testStr)).toBeTruthy();
  });

  var LOWERCASE_REGEX = new RegExp(/[a-z]+/);
  var UPPERCASE_REGEX = new RegExp(/[A-Z]+/);

  it('full lowercase matches regex', function() {
    var testStr = generatorService.FULL_LOWERCASE;
    expect(LOWERCASE_REGEX.test(testStr)).toBeTruthy();
    expect(UPPERCASE_REGEX.test(testStr)).toBeFalsy();
  });

  it('full uppercase set', function() {
    var testStr = generatorService.FULL_UPPERCASE;
    expect(angular.isDefined(testStr)).toBeTruthy();
    expect(angular.isString(testStr)).toBeTruthy();
  });

  it('full uppercase matches regex', function() {
    var testStr = generatorService.FULL_UPPERCASE;
    expect(LOWERCASE_REGEX.test(testStr)).toBeFalsy();
    expect(UPPERCASE_REGEX.test(testStr)).toBeTruthy();
  });

  it('non similar lowercase set', function() {
    var testStr = generatorService.NON_SIMILAR_LOWERCASE;
    expect(angular.isDefined(testStr)).toBeTruthy();
    expect(angular.isString(testStr)).toBeTruthy();
  });

  it('non similar lowercase missing characters', function() {
    var testStr = generatorService.NON_SIMILAR_LOWERCASE;
    for (var i in ['i', 'o', 'l']) {
      expect(testStr.indexOf(i)).toEqual(-1);
    }
  });

  it('non similar uppercase missing characters', function() {
    var testStr = generatorService.NON_SIMILAR_UPPERCASE;
    for (var i in ['I', 'O']) {
      expect(testStr.indexOf(i)).toEqual(-1);
    }
    expect(testStr.indexOf('L')).toBeGreaterThan(-1);
  });
});
