export default {
  validators() {
    return {
      required: value => !!value || "This field is required",
      password: value => value.length >= 6 || "Less than 6 characters entered",
      equel(v) {
        return value => v == value || "This passwords do not match";
      },
      email: value =>
        /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/.test(
          value
        ),
    };
  },
  elValidators() {
    return {
      required: true,
      message: 'Please pick a date',
    }
  },
  isValidLink(link) {
    if (!link) {
      return false; // Empty string is not a valid link
    }

    // Check for HTTPS link
    const httpsRegex = /^(https?:\/\/)([\w.-]+)\.([a-z]{2,6}\.?)(\/[\w.-]*)*\/?$/i;

    // Check for Telegram link
    const telegramRegex = /^(https?:\/\/)?(t\.me\/|telegram\.me\/)([\w\d]+)\/?$/i; // Handles both with and without https

    return httpsRegex.test(link) || telegramRegex.test(link);
  }
};
