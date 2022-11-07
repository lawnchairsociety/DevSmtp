namespace DevSmtp.Core.Utilities
{
    public static class StringExtensions
    {
        public static bool IsEmpty(this string? candidate) => string.IsNullOrWhiteSpace(candidate);
    }
}
