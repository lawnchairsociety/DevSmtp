namespace DevSmtp.Core.Commands
{
    public class SomlException : Exception
    {
        public SomlException(string message)
            : base(message)
        {
        }

        public SomlException(string message, Exception innerException)
            : base(message, innerException)
        {
        }
    }
}
