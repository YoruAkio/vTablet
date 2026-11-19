import 'package:flutter/material.dart';
import 'package:vtablet/configs.dart';
import 'package:vtablet/services/states.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_gen/gen_l10n/localizations.dart';

import 'home.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await ConfigManager.initManager();
  runApp(const App());
}

class App extends StatefulWidget {
  const App({super.key});

  static void setLocale(BuildContext context, Locale? newLocale) {
    _AppState? state = context.findAncestorStateOfType<_AppState>();
    state?.setLocale(newLocale);
  }

  @override
  State<App> createState() => _AppState();
}

class _AppState extends State<App> {
  bool useMaterial3 = true;
  ThemeMode themeMode = ThemeMode.system;
  Locale? _locale;

  @override
  void initState() {
    super.initState();
    String savedLang = Configs.appLanguage.get();
    if (savedLang == 'en') {
      _locale = const Locale('en');
    } else if (savedLang == 'id') {
      _locale = const Locale('id');
    } else if (savedLang == 'zh') {
      _locale = const Locale('zh');
    } else {
      _locale = null; // Follow system, default to English fallback
    }
  }

  void setLocale(Locale? newLocale) {
    setState(() {
      _locale = newLocale;
    });
  }

  bool get useLightMode {
    switch (themeMode) {
      case ThemeMode.system:
        return WidgetsBinding.instance.platformDispatcher.platformBrightness ==
            Brightness.light;
      case ThemeMode.light:
        return true;
      case ThemeMode.dark:
        return false;
    }
  }

  void handleBrightnessChange() {
    setState(() {
      themeMode = useLightMode ? ThemeMode.dark : ThemeMode.light;
    });
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      title: 'vTablet',
      locale: _locale,
      localizationsDelegates: const [
        AppLocalizations.delegate,
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
      ],
      supportedLocales: const [
        Locale('en'), // English (Default)
        Locale('id'), // Indonesian
        Locale('zh'), // Chinese
      ],
      localeResolutionCallback: (deviceLocale, supportedLocales) {
        if (_locale != null) {
          return _locale;
        }
        if (deviceLocale != null) {
          for (var supported in supportedLocales) {
            if (supported.languageCode == deviceLocale.languageCode) {
              return supported;
            }
          }
        }
        return const Locale('en'); // Default to English
      },
      themeMode: themeMode,
      theme: ThemeData(
        useMaterial3: useMaterial3,
        brightness: Brightness.light,
      ),
      darkTheme: ThemeData(
        useMaterial3: useMaterial3,
        brightness: Brightness.dark,
      ),
      home: Home(
        title: "vTablet",
        handleBrightnessChange: handleBrightnessChange,
      ),
    );
  }
}
