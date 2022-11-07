namespace DevSmtp.Core.Commands
{
    public class HelpException : Exception
    {
        public HelpException(string message)
            : base(message)
        {
        }

        public HelpException(string message, Exception innerException)
            : base(message, innerException)
        {
        }
    }
}
