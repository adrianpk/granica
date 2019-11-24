# Changelog

## <a name="20191124"></a>20191124 - Embedded translations

Localization files now are mbedded in executable as any other resource.

## <a name="20191123"></a>20191123 - Internationalization

Not fully implemented but basically this is how it works.

The I18N middleware currently tries to read the `lang` field submitted in POST, PUT, PATCH request; alternatively that same field value but read from the query string is used for all methods.

If still not found the value to be used to decide the language to be use will be associated to the one reported by `Accept-Language` request HTTP header.

In any other case the default language (English) will be used.

Soon, priority to decide the language will be the lang user chooses at sign-up time and/or another set by user using a gui switch.

A general clean-up is still needed, eventually a large part of the implementation will be moved to web module.

Finally, as with other assets, localization files under `assets/web/embed/i18n` will be embedded directly into the executable. Temporarily at this stage of development they are accessed from the filesytem.
