import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:mobile/pages/home_page/home_page.dart';
import 'package:mobile/pages/login_page/login_page.dart';

final GoRouter router = GoRouter(
  routes: <RouteBase>[
    GoRoute(
      path: '/',
      builder: (BuildContext context,GoRouterState state){
        return const LoginPage();
      },
    ),
    GoRoute(  
        path: '/',
        builder: (BuildContext context,GoRouterState state){
          return const LoginPage();
        },
    ),
  ]
);