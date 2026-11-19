import 'package:flutter/material.dart';
import 'package:flutter_gen/gen_l10n/localizations.dart';
import 'package:vtablet/configs.dart';
import 'package:vtablet/main.dart';

class AppearancePage extends StatefulWidget {
  const AppearancePage({
    super.key,
  });

  @override
  State<AppearancePage> createState() => _AppearancePageState();
}

class _AppearancePageState extends State<AppearancePage> {
  String _selectedLang = Configs.appLanguage.get();

  void _onLangChanged(String? lang) {
    if (lang == null) return;
    setState(() {
      _selectedLang = lang;
      Configs.appLanguage.set(lang);
    });

    Locale? newLocale;
    if (lang == 'en') {
      newLocale = const Locale('en');
    } else if (lang == 'id') {
      newLocale = const Locale('id');
    } else if (lang == 'zh') {
      newLocale = const Locale('zh');
    } else if (lang == 'ru') {
      newLocale = const Locale('ru');
    } else if (lang == 'pt') {
      newLocale = const Locale('pt');
    } else if (lang == 'ja') {
      newLocale = const Locale('ja');
    } else if (lang == 'fr') {
      newLocale = const Locale('fr');
    } else if (lang == 'de') {
      newLocale = const Locale('de');
    } else if (lang == 'pl') {
      newLocale = const Locale('pl');
    } else {
      newLocale = null;
    }
    App.setLocale(context, newLocale);
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    return ListView(
      padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 16.0),
      children: [
        Card.outlined(
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    const Icon(Icons.language, size: 24),
                    const SizedBox(width: 12),
                    Text(
                      l10n.appearanceLanguageTitle,
                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                            fontWeight: FontWeight.bold,
                          ),
                    ),
                  ],
                ),
                const SizedBox(height: 12),
                RadioListTile<String>(
                  title: Text(l10n.appearanceLanguageSystem),
                  subtitle: const Text('English (Default)'),
                  value: 'system',
                  groupValue: _selectedLang,
                  onChanged: _onLangChanged,
                ),
                RadioListTile<String>(
                  title: const Text('English'),
                  subtitle: const Text('English'),
                  value: 'en',
                  groupValue: _selectedLang,
                  onChanged: _onLangChanged,
                ),
                RadioListTile<String>(
                  title: const Text('Bahasa Indonesia'),
                  subtitle: const Text('Indonesian'),
                  value: 'id',
                  groupValue: _selectedLang,
                  onChanged: _onLangChanged,
                ),
                RadioListTile<String>(
                  title: const Text('中文'),
                  subtitle: const Text('Chinese'),
                  value: 'zh',
                  groupValue: _selectedLang,
                  onChanged: _onLangChanged,
                ),
                RadioListTile<String>(
                  title: const Text('Русский'),
                  subtitle: const Text('Russian'),
                  value: 'ru',
                  groupValue: _selectedLang,
                  onChanged: _onLangChanged,
                ),
                RadioListTile<String>(
                  title: const Text('Português (Brasil)'),
                  subtitle: const Text('Portuguese'),
                  value: 'pt',
                  groupValue: _selectedLang,
                  onChanged: _onLangChanged,
                ),
                RadioListTile<String>(
                  title: const Text('日本語'),
                  subtitle: const Text('Japanese'),
                  value: 'ja',
                  groupValue: _selectedLang,
                  onChanged: _onLangChanged,
                ),
                RadioListTile<String>(
                  title: const Text('Français'),
                  subtitle: const Text('French'),
                  value: 'fr',
                  groupValue: _selectedLang,
                  onChanged: _onLangChanged,
                ),
                RadioListTile<String>(
                  title: const Text('Deutsch'),
                  subtitle: const Text('German'),
                  value: 'de',
                  groupValue: _selectedLang,
                  onChanged: _onLangChanged,
                ),
                RadioListTile<String>(
                  title: const Text('Polski'),
                  subtitle: const Text('Polish'),
                  value: 'pl',
                  groupValue: _selectedLang,
                  onChanged: _onLangChanged,
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }
}

class AppearanceDestination extends StatelessWidget {
  const AppearanceDestination({
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return NavigationDestination(
      icon: const Icon(Icons.color_lens_outlined),
      selectedIcon: const Icon(Icons.color_lens),
      label: AppLocalizations.of(context)!.appearancePageTitle,
    );
  }
}

NavigationRailDestination appearanceRailDestination(BuildContext context) {
  return NavigationRailDestination(
    icon: const Icon(Icons.color_lens_outlined),
    selectedIcon: const Icon(Icons.color_lens),
    label: Text(AppLocalizations.of(context)!.appearancePageTitle),
  );
}
