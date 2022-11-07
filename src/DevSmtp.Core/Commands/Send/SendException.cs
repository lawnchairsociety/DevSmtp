namespace DevSmtp.Core.Commands
{
    public class SendException : Exception
    {
        public SendException(string message)
            : base(message)
        {
        }

        public SendException(string message, Exception innerException)
            : base(message, innerException)
        {
        }
    }
}
