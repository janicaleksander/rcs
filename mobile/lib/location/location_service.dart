import 'dart:io';

import 'package:location/location.dart';

class LocationService {
  final Location _location = Location();
  LocationData? _locationData;

  LocationService();

  Future<bool> isPermission() async {
    var serviceEnabled = await _location.serviceEnabled();
    if (!serviceEnabled) {
      serviceEnabled = await _location.requestService();
      if (!serviceEnabled) {
        return false;
      }
    }
    var permissionStatus = await _location.hasPermission();
    if (permissionStatus == PermissionStatus.denied) {
      permissionStatus = await _location.requestPermission();
      if (permissionStatus != PermissionStatus.granted) {
        return false;
      }
    }
    /*
    _location.changeSettings(
      accuracy: LocationAccuracy.high,
      interval: 5000,
      distanceFilter: 3,
    );
     */
    _location.changeSettings(
      accuracy: LocationAccuracy.high,
      interval: 1000,
      distanceFilter: 0,
    );

    return true;
  }

  Future<void> updateLocation() async {
    _locationData = await _location.getLocation();
  }

  LocationData? getLocationData() {
    return _locationData;
  }
    Future<bool> enableBackgroundMode() async {
      bool _bgModeEnabled = await _location.isBackgroundModeEnabled();
      if (_bgModeEnabled) {
        return true;
      } else {
        try {
          await _location.enableBackgroundMode();
        } catch (e) {
          print(e.toString());
        }
        try {
          _bgModeEnabled = await _location.enableBackgroundMode();
        } catch (e) {
          print(e.toString());
        }
        print(_bgModeEnabled); //True!
        return _bgModeEnabled;
      }
    }
  void registerHandler(Future<void> Function() handler) {
    handler();
    _location.onLocationChanged.listen((LocationData currentLocation) {
      _locationData = currentLocation;
      handler();
    });
  }
}
