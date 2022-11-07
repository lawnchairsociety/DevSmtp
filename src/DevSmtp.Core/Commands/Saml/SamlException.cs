namespace DevSmtp.Core.Commands
{
    public class SamlException : Exception
    {
        public SamlException(string message)
            : base(message)
        {
        }

        public SamlException(string message, Exception innerException)
            : base(message, innerException)
        {
        }
    }
}
